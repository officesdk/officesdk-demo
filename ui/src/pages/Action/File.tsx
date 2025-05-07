import { useEffect, useState, RefCallback } from "react"
import api from '@/services/case'
import { useParams, useLocation } from 'umi'
import { message, Button } from "antd"
import { createSDK } from '@officesdk/web';
import { getFileTypeFromExt } from '@/utils/common'
import { useIntl } from 'umi';
import { SaveOutlined } from "@ant-design/icons"

export default function File() {
  const intl = useIntl();
  const { id } = useParams()
  const { pathname } = useLocation()
  const [endpoint, setEndpoint] = useState("")
  const [fileExt, setFileExt] = useState("")
  const [editor, setEditor] = useState<any>(null)
  const isEdit = pathname.includes("collab")
  const [loading, setLoading] = useState(false)
  const [token, setToken] = useState("")

  useEffect(() => {
    getFilePreview()
  }, [])

  useEffect(() => {
    (window as any).editor = editor
  }, [editor])

  // 获取预览html文件
  function getFilePreview() {
    if (!id) {
      message.warning("404")
      return
    }
    api.GetCasePreview(id?.toString()).then((res) => {
      setEndpoint(res.endpoint)
      setToken(res.token)
      setFileExt(res.file.ext)
    })
  }

  const loadSDK: RefCallback<HTMLElement> = async (el) => {
    if (el && endpoint && !editor) {
      try {
        const sdk = createSDK({
          // 服务器地址
          endpoint: endpoint,
          // 路径,暂时不用
          // path: 'v1/api/file/page',
          fileId: id,
          mode: isEdit ? 'standard' : 'preview',
          role: isEdit ? 'editor' : 'viewer',
          // 当前语言
          lang: intl.locale === 'en-US' ? 'en-US' : 'zh-CN',
          // 装在iframe的根结点
          root: el,
          fileType: getFileTypeFromExt(fileExt),
          token: token,
          settings: {
            menu: {
              custom: [
                {
                  name: 'test',
                  type: 'button',
                  label: '按钮测试',
                  callback: (): void => {
                    console.log('按钮测试');
                  },
                },
              ],
            },
          },
          userQuery: {
            userName: "demo",
            userId: "1"
          }
        });
        const editor = await sdk.connect()
        setEditor(editor)
        console.info(editor)
      } catch (error) {
        console.error('Failed to create SDK:', error);
      }
    }
  };

  return (
    <div style={{ width: "100%", height: "100vh" }}>
      {
        isEdit && (
          <div style={{ display: "flex", justifyContent: "flex-end", padding: "5px", borderBottom: "1px solid #eee" }}>
            <Button loading={loading} style={{ width: "100px" }}
              icon={<SaveOutlined />}
              onClick={() => {
                setLoading(true)
                editor?.content.save().then((res: any) => {
                  setLoading(false)
                  message.success(intl.formatMessage({ id: 'action.save.success' }))
                })
              }}>
              {intl.formatMessage({ id: 'action.save' })}
            </Button>
          </div>
        )
      }
      <div style={{ width: "100%", height: isEdit ? "calc(100vh - 48px)" : "100vh" }} ref={loadSDK}>
      </div>
    </div>
  )
}