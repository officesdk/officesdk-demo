import { Button, message, Table, Upload, Popconfirm, Tooltip } from 'antd'
import type { TableProps, UploadProps } from 'antd'
import { useEffect, useState } from 'react'
import api from '@/services/action'
import { Link } from 'react-router-dom'

export default function Action() {
  const [data, setData] = useState([])
  const [loading, setLoading] = useState(false)
  const [listLoading, setListLoading] = useState(false)
  const columns: TableProps['columns'] = [
    {
      title: '',
      dataIndex: 'index',
      key: 'index',
      width: 60
    },
    {
      title: '文件名',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '创建时间',
      dataIndex: 'create_time',
      key: 'create_time',
      width: 300,
      render: (text) => {
        return <span>{transDate(text)}</span>
      }
    },
    {
      title: '',
      dataIndex: 'action',
      key: 'action',
      width: 100,
      render: (_, record) => {
        return <>
          <Link to={`/file/${record.id}`} target="_blank">预览</Link>
          {
            !record.id.includes("case_") && (
              <Popconfirm
                title="删除"
                description="确认删除？"
                onConfirm={() => handleRemove(record.id)}
                okText="确认"
                cancelText="取消"
              >
                <a style={{ marginLeft: "8px" }}>删除</a>
              </Popconfirm>
            )
          }
        </>
      }
    },
  ]

  useEffect(() => {
    getFileList()
  }, [])

  const uploadProps: UploadProps = {
    name: 'file',
    customRequest(option: any) {
      try {
        setLoading(true)
        api.UploadFile(option.file).then(res => {
          setLoading(false)
          message.success("上传成功")
          getFileList()
          option.onSuccess(res.data.filePath)
        }).catch(err => {
          setLoading(false)
          option.onError(err?.data?.message)
        })
      } catch (error) {
        option.onError(error)
        setLoading(false)
      }
    },
    showUploadList: false,
  }
  // 获取文件列表
  const getFileList = () => {
    setListLoading(true)
    api.GetFiles().then(res => {
      setData(res?.sort((a: API.FileData, b: API.FileData) => b.create_time - a.create_time)?.map((item: API.FileData, index: number) => {
        return {
          index: index + 1,
          ...item
        }
      }))
    }).finally(() => {
      setListLoading(false)
    })
  }

  const transDate = (date: number) => {
    if (date == 0) return ""
    const d = new Date(date * 1000)
    return d.toLocaleString()
  }

  const handleRemove = (fileId: string) => {
    api.DeleteFiles(fileId).then(res => {
      setData((origin) => {
        return origin.filter((item: API.FileData) => item.id !== fileId)
      })
    })
  }
  return <>
    {/* {
      data?.length >= 20 ? (
        <Tooltip title="文件列表限制数量20">
          <Button disabled type='primary' style={{ marginBottom: "10px" }}>上传文件</Button>
        </Tooltip>
      ) : ( */}
    <Upload {...uploadProps}>
      <Button type='primary' loading={loading} style={{ marginBottom: "10px" }}>上传文件</Button>
    </Upload>
    {/* )
    } */}
    <Table loading={listLoading} columns={columns} dataSource={data} pagination={false} rowKey="id"></Table>
  </>
}