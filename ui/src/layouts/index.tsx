import { Outlet, useLocation } from 'umi';
import { Layout, Menu, ConfigProvider } from 'antd';
import type { MenuProps } from 'antd';
import styles from './index.less'
import { Link } from 'react-router-dom'
import zhCN from 'antd/es/locale/zh_CN';
export default function Layouts() {
  const { Header, Content, Sider } = Layout;
  let currentRout = useLocation().pathname;
  // 右侧菜单项
  const menuItems: MenuProps['items'] = [
    {
      label: <Link to={`/showcase`}>示例</Link>,
      key: "showcase"
    },
    {
      label: <Link to={`/action`}>功能</Link>,
      key: "action"
    },
  ]

  // 获取默认高亮菜单
  function getDefaultMenu() {
    let menu = currentRout.replace("/", "")
    return menu || "showcase"
  }

  return (
    <ConfigProvider locale={zhCN}>
      <Layout className={styles.Layout}>
        <Header className={styles.Header}>
          <div className={styles.Introduce}>
            <span className={styles.Title}>极速版 SDK 测试版本</span>
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
