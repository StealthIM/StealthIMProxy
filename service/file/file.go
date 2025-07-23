package file

import (
	pb "StealthIMProxy/StealthIM.FileAPI"
	"StealthIMProxy/conns"
	"StealthIMProxy/errorcode"
	"context"
	"encoding/binary"
	"encoding/json" // 导入io包
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func call() (*pb.StealthIMFileAPIClient, error) {
	conn, err := conns.PoolFileAPI.ChooseConn()
	if err != nil {
		return nil, err
	}
	cli := pb.NewStealthIMFileAPIClient(conn)
	return &cli, nil
}

func upload(c *gin.Context) {
	uid := int32(c.GetInt("uid"))
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	pbconn, err := call()
	if err != nil {
		conn.WriteJSON(gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  "Server internal network error",
			},
		})
		return
	}
	stream, err := (*pbconn).Upload(c)
	defer stream.CloseSend()
	// metadata
	msgType, msg, err := conn.ReadMessage()
	if err != nil {
		return
	}
	if msgType != websocket.TextMessage {
		conn.WriteJSON(gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  "Metadata must be text message",
			},
			"type": "metadata",
		})
		return
	}
	metadata := &meatdataRequest{}
	err = json.Unmarshal(msg, metadata)
	if err != nil {
		conn.WriteJSON(gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  "Metadata format error",
			},
			"type": "metadata",
		})
		return
	}
	if metadata.Size == "" || metadata.Size == "0" || metadata.Groupid == "" || metadata.Groupid == "0" || metadata.Hash == "" {
		conn.WriteJSON(gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  "Metadata format error",
			},
			"type": "metadata",
		})
		return
	}
	grpid, err := strconv.ParseInt(metadata.Groupid, 10, 64)
	if err != nil {
		conn.WriteJSON(gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  "Metadata format error",
			},
			"type": "metadata",
		})
		return
	}
	filesize, err := strconv.ParseInt(metadata.Size, 10, 64)
	if err != nil {
		conn.WriteJSON(gin.H{
			"result": gin.H{
				"code": errorcode.ProxyBadRequest,
				"msg":  "Metadata format error",
			},
			"type": "metadata",
		})
		return
	}
	err = stream.Send(&pb.UploadRequest{Data: &pb.UploadRequest_Metadata{
		Metadata: &pb.Upload_FileMetaData{
			Totalsize:     filesize,
			UploadGroupid: grpid,
			UploadUid:     uid,
			Hash:          metadata.Hash,
		},
	}})
	if err != nil {
		conn.WriteJSON(gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  "Server internal network error",
			},
			"type": "metadata",
		})
		return
	}
	in, err := stream.Recv()
	if err != nil {
		conn.WriteJSON(gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  "Network error",
			},
			"type": "metadata",
		})
		return
	}
	if in.Result.Code != errorcode.Success {
		conn.WriteJSON(gin.H{
			"result": gin.H{
				"code": in.Result.Code,
				"msg":  in.Result.Msg,
			},
		})
		return
	}
	var metaResp *pb.Upload_MetaResponse
	metaResp = in.GetMeta()
	if metaResp == nil {
		conn.WriteJSON(gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  "Network error",
			},
			"type": "metadata",
		})
		return
	}
	conn.WriteJSON(gin.H{
		"result": gin.H{
			"code": errorcode.Success,
			"msg":  "",
		},
		"type": "metadata",
	})

	go func() {
		for {
			select {
			case <-c.Done():
				return
			default:
				in, err := stream.Recv()
				if err != nil {
					conn.WriteJSON(gin.H{
						"result": gin.H{
							"code": errorcode.ServerInternalNetworkError,
							"msg":  "Network error",
						},
						"type": "complete",
					})
					conn.Close()
					return
				}
				if blkResp := in.GetMeta(); blkResp != nil {
					conn.WriteJSON(gin.H{
						"result": gin.H{
							"code": in.Result.Code,
							"msg":  in.Result.Msg,
						},
						"type": "metadata",
					})
					conn.Close()
					return
				} else if blkResp := in.GetBlock(); blkResp != nil {
					conn.WriteJSON(gin.H{
						"result": gin.H{
							"code": in.Result.Code,
							"msg":  in.Result.Msg,
						},
						"type":    "block",
						"blockid": blkResp.Blockid,
					})
				} else if clpResp := in.GetComplete(); clpResp != nil {
					conn.WriteJSON(gin.H{
						"result": gin.H{
							"code": in.Result.Code,
							"msg":  in.Result.Msg,
						},
						"type": "complete",
					})
					conn.Close()
					return
				}
			}
		}
	}()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if msgType != websocket.BinaryMessage {
			conn.WriteJSON(gin.H{
				"result": gin.H{
					"code": errorcode.ProxyBadRequest,
					"msg":  "Data format error",
				},
				"type": "complete",
			})
			return
		}
		if len(msg) < 5 {
			conn.WriteJSON(gin.H{
				"result": gin.H{
					"code": errorcode.ProxyBadRequest,
					"msg":  "Data length error",
				},
				"type": "complete",
			})
			return
		}
		blockid := int32(binary.LittleEndian.Uint32(msg[:4]))
		dataBody := msg[4:]
		err = stream.Send(&pb.UploadRequest{Data: &pb.UploadRequest_File{
			File: &pb.Upload_FileBlock{
				Blockid: blockid,
				File:    dataBody,
			},
		}})
		if err != nil {
			conn.WriteJSON(gin.H{
				"result": gin.H{
					"code": errorcode.ServerInternalNetworkError,
					"msg":  "Save data error",
				},
				"type":    "block",
				"blockid": int32(blockid),
			})
		}
	}
}

func getInfo(c *gin.Context) {
	var fileHash string = c.Param("filehash")

	cli, err := call()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ret, err := (*cli).GetFileInfo(ctx, &pb.GetFileInfoRequest{
		Hash: fileHash,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  "Network error",
			},
		})
		return
	}
	if ret.Result.Code != errorcode.Success {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": ret.Result.Code,
				"msg":  ret.Result.Msg,
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"code": errorcode.Success,
			"msg":  "",
		},
		"size": fmt.Sprintf("%d", ret.Size),
	})

}

func download(c *gin.Context) {
	var fileHash string = c.Param("filehash")
	rangeHeader := c.GetHeader("Range")

	var start, end int64 = 0, 0
	if rangeHeader != "" {
		// 格式: bytes=start-end
		rangeParts := rangeHeader[len("bytes="):]
		dashIndex := -1
		for i, r := range rangeParts {
			if r == '-' {
				dashIndex = i
				break
			}
		}

		if dashIndex != -1 {
			startStr := rangeParts[:dashIndex]
			endStr := rangeParts[dashIndex+1:]

			if startStr != "" {
				s, err := strconv.ParseInt(startStr, 10, 64)
				if err == nil {
					start = s
				}
			}
			if endStr != "" {
				e, err := strconv.ParseInt(endStr, 10, 64)
				if err == nil {
					end = e
				}
			}
		}
	}

	cli, err := call()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  err.Error(),
			},
		})
		return
	}

	// 获取文件信息以获取文件总大小
	fileInfoRet, err := (*cli).GetFileInfo(context.Background(), &pb.GetFileInfoRequest{
		Hash: fileHash,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  "Network error",
			},
		})
		return
	}
	if fileInfoRet.Result.Code != errorcode.Success {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": fileInfoRet.Result.Code,
				"msg":  fileInfoRet.Result.Msg,
			},
		})
		return
	}
	totalSize := int64(fileInfoRet.Size) // 将totalSize转换为int64

	// 如果end为0，表示请求到文件末尾
	if end == 0 {
		end = totalSize - 1
	}

	// 确保start和end在有效范围内
	if start < 0 {
		start = 0
	}
	if end >= totalSize {
		end = totalSize - 1
	}
	if start > end {
		start = 0
		end = totalSize - 1
	}

	c.Header("Content-Type", "application/stim-partid-binary")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10)+"/"+strconv.FormatInt(totalSize, 10))
	c.Writer.WriteHeader(http.StatusPartialContent)
	c.Writer.Flush()

	stream, err := (*cli).Download(context.Background(), &pb.DownloadRequest{
		Hash:  fileHash,
		Start: (start),
		End:   (end + 1), // 原API是包含头不包含尾的，所以这里需要加1
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"code": errorcode.ServerInternalNetworkError,
				"msg":  "Network error",
			},
		})
		return
	}

	c.Stream(func(w io.Writer) bool {
		for {
			resp, err := stream.Recv()
			if err != nil {
				return false
			}
			if data := resp.GetFile(); data != nil {
				sendDataHeader := make([]byte, 4)
				binary.LittleEndian.PutUint32(sendDataHeader, uint32(data.Blockid))
				sendDataLen := make([]byte, 4)
				binary.LittleEndian.PutUint32(sendDataLen, uint32(len(data.File)))
				sendData := append(sendDataHeader, sendDataLen...)
				sendData = append(sendData, data.File...)
				if _, err := w.Write(sendData); err != nil {
					return false
				}
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			} else if data := resp.GetResult(); data != nil {
				sendDataHeader := []byte{0xff, 0xff, 0xff, 0xff}
				jsonObj := gin.H{
					"result": gin.H{
						"code": data.Code,
						"msg":  data.Msg,
					},
				}
				jsonBytes, err := json.Marshal(jsonObj)
				if err != nil {
					return false
				}
				sendDataLen := make([]byte, 4)
				binary.LittleEndian.PutUint32(sendDataLen, uint32(len(jsonBytes)))
				sendData := append(sendDataHeader, sendDataLen...)
				sendData = append(sendData, jsonBytes...)
				if _, err := w.Write(sendData); err != nil {
					return false
				}
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
				return false
			}
		}
	})
}
