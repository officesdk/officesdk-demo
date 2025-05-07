import { Outlet, useLocation, useIntl, setLocale } from 'umi';
import { useEffect } from 'react'
import { Layout, Menu, ConfigProvider, Popover, Button } from 'antd';
import type { MenuProps } from 'antd';
import styles from './index.less'
import { Link } from 'react-router-dom'
import zhCN from 'antd/es/locale/zh_CN';
import enUS from 'antd/es/locale/en_US';

export default function Layouts() {
  const { Header, Content, Sider } = Layout;
  let currentRout = useLocation().pathname;
  const intl = useIntl();
  // 右侧菜单项
  const menuItems: MenuProps['items'] = [
    {
      label: <Link to={`/showcase`}>{intl.formatMessage({ id: 'action.list.showcase' })}</Link>,
      key: "showcase"
    },
    {
      label: <Link to={`/action`}>{intl.formatMessage({ id: 'action.list.action' })}</Link>,
      key: "action"
    },
  ]

  useEffect(() => {
    // 初始化语言
    if (localStorage.getItem('Locale')) {
      setLocale(localStorage.getItem('Locale') || navigator.language, false)
    } else {
      setLocale(navigator.language == "en" ? "en-US" : "zh-CN", false)
      localStorage.setItem('Locale', navigator.language == "en" ? "en-US" : "zh-CN")
    }
  }, [])

  // 获取默认高亮菜单
  function getDefaultMenu() {
    let menu = currentRout.replace("/", "")
    return menu || "showcase"
  }

  // 更换语言
  const changeLang = (lang: string) => {
    if (lang == "zh-CN") {
      setLocale("en-US", false)
      localStorage.setItem('Locale', "en-US")
    } else if (lang == "en-US") {
      setLocale("zh-CN", false)
      localStorage.setItem('Locale', "zh-CN")
    }
  }

  return (
    <ConfigProvider locale={intl.locale == "en-US" ? enUS : zhCN}>
      <Layout className={styles.Layout}>
        <Header className={styles.Header}>
          <div className={styles.Introduce}>
            <span className={styles.Title}>{intl.formatMessage({ id: 'title' })}</span>
            <Popover content={intl.locale == "en-US" ? "English / 中文" : "中文 / English"} placement="bottomRight">
              <Button type='link' className={styles.Lang} onClick={() => changeLang(intl.locale)}>
                {intl.locale == "en-US" ? "EN" : "中文"}
              </Button>
            </Popover>
          </div>
        </Header>
        <Layout>
          <Sider className={styles.Sider}>
            <Menu items={menuItems} defaultSelectedKeys={[getDefaultMenu()]} />
          </Sider>
          <Content className={styles.Content}>
            <div className={styles.Body} ><Outlet /> </div>
          </Content>
        </Layout>
      </Layout>
    </ConfigProvider>
  );
}
