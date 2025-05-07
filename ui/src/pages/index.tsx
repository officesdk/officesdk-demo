import { Button, Card } from 'antd';
import styles from './index.less'
import { getImage } from '@/utils/common'
import { exampleCase, exampleCaseEn } from '@/constants/case'
import { useIntl } from 'umi';
import { useState, useEffect } from 'react';

export default function HomePage() {
    const intl = useIntl();
    const { Meta } = Card;
    const [caseList, setCaseList] = useState(intl.locale == 'zh-CN' ? exampleCase : exampleCaseEn)

    useEffect(() => {
        setCaseList(intl.locale == 'zh-CN' ? exampleCase : exampleCaseEn)
    }, [intl.locale])

    return (
        <div className={styles.Content}>
            {caseList && caseList.map((item, index) => {
                return (
                    <Card
                        key={index}
                        hoverable
                        className={styles.exampleCards}
                        cover={<span>{item.title}</span>}
                    >
                        <Meta description={
                            <div className={styles.exampleBody}>
                                <img className={styles.imgCase} alt="example" src={getImage(item.type, intl.locale)?.image} />
                                <div className={styles.hoverCase} style={{ padding: intl.locale == 'zh-CN' ? "" : "50px" }}>
                                    <Button
                                        disabled={!item.support}
                                        onClick={() => window.open(`/file/${item.id}`)}
                                    >
                                        {intl.formatMessage({ id: 'showcase.to' }, { text: getImage(item.type, intl.locale)?.text })}
                                    </Button>
                                    {item.editable && (
                                        <Button
                                            disabled={!item.support}
                                            onClick={() => window.open(`/collab/${item.id}`)}
                                        >
                                            {intl.formatMessage({ id: 'showcase.to.edit' }, { text: getImage(item.type, intl.locale)?.text })}
                                        </Button>
                                    )}
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
