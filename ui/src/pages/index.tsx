import { Button, Card } from 'antd';
import styles from './index.less'
import { getImage } from '@/utils/common'
import { exampleCase } from '@/constants/case'

export default function HomePage() {
    const { Meta } = Card;
    return (
        <div className={styles.Content}>
            {exampleCase && exampleCase.map((item, index) => {
                return (
                    <Card
                        key={index}
                        hoverable
                        className={styles.exampleCards}
                        cover={<span>{item.title}</span>}
                    >
                        <Meta description={
                            <div className={styles.exampleBody}>
                                <img className={styles.imgCase} alt="example" src={getImage(item.type)?.image} />
                                <div className={styles.hoverCase}>
                                    <Button
                                        disabled={!item.support}
                                        onClick={() => window.open(`/case?type=${item.type}`)}
                                    >
                                        跳转{getImage(item.type)?.text}预览
                                    </Button>
                                </div>
                            </div>
                        }
                        />
                    </Card>
                )
            })}
        </div>
    );
}
