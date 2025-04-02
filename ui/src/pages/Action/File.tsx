import { useEffect, useState } from "react"
import api from '@/services/case'
import { useParams } from 'umi'
import { message } from "antd"

export default function File() {
  const { id } = useParams()
  // 预览结果文件
  const [fileUrl, setFileUrl] = useState("")
  useEffect(() => {
    getFilePreview()
  }, [])
  // 获取预览html文件
  function getFilePreview() {
    if (!id) {
      message.warning("404")
      return
    }
    api.GetCasePreview(id?.toString()).then((res) => {
      setFileUrl(res.url)
    })
  }

  return <>
    <iframe
      allow="fullscreen"
      allowFullScreen={true}
      style={{ border: 0, width: "100%", height: "100vh", display: "block" }}
      src={fileUrl}
    >
      {fileUrl}
    </iframe>
  </>
}