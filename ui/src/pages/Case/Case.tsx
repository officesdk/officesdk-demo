import { useEffect, useState } from "react"
import api from '@/services/case'
import { Result, Button } from 'antd'
import { history } from 'umi'

export default function Case() {
  // 预览加载状态
  const [status, setStatus] = useState(0)
  // 预览结果文件
  const [fileUrl, setFileUrl] = useState("")
  useEffect(() => {
    getFilePreview()
  }, [])
  // 获取预览html文件
  function getFilePreview() {
    setStatus(1)
    const queryParams = new URLSearchParams(location.search)
    let type = queryParams.get("type")
    let caseItem = {} as { name: string, url: string }
    switch (type) {
      case "image":
        caseItem = { name: "image", url: 'image.png' }
        break;
      case "word_zh":
        caseItem = { name: "word_zh", url: 'word_zh.docx' }
        break;
      case "word_en":
        caseItem = { name: "word_en", url: 'word_en.docx' }
        break;
      case "excel_zh":
        caseItem = { name: "excel_zh", url: 'excel_zh.xlsx' }
        break;
      case "excel_en":
        caseItem = { name: "excel_en", url: 'excel_en.xlsx' }
        break;
      case "ppt_zh":
        caseItem = { name: "ppt_zh", url: 'ppt_zh.pptx' }
        break;
      case "ppt_en":
        caseItem = { name: "ppt_en", url: 'ppt_en.pptx' }
        break;
      default:
        setStatus(2)
        return
    }
    api.GetCasePreview(`case_${caseItem.url}`).then((res) => {
      setFileUrl(res.url)
    })
  }

  return <>
    {status == 2 ?
      <Result
        status="404"
        title="404"
        subTitle="Sorry, the type you want is not supported."
        extra={<Button type="primary" onClick={() => history.push('/')}>Back Home</Button>}
      />
      :
      <iframe
        allow="fullscreen"
        allowFullScreen={true}
        style={{ border: 0, width: "100%", height: "100vh", display: "block" }}
        src={fileUrl}
      >
        {fileUrl}
      </iframe>
    }
  </>
}