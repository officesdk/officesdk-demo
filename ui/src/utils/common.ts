
import word_zh from '../assets/word_zh.png'
import word_en from '../assets/word_en.png'
import excel_zh from '../assets/excel_zh.png'
import excel_en from '../assets/excel_en.png'
import ppt_zh from '../assets/ppt_zh.png'
import ppt_en from '../assets/ppt_en.png'
import image from '../assets/image.png'

// 获取对应case信息
export function getImage(type: string, lang: string = 'zh-CN'): { image: string, text: string } {
  switch (type) {
    case 'excel_zh':
      return { image: excel_zh, text: '表格' }
    case 'excel_en':
      return { image: excel_en, text: 'Excel' }
    case 'ppt_zh':
      return { image: ppt_zh, text: '演示文档' }
    case 'ppt_en':
      return { image: ppt_en, text: 'PPT' }
    case 'word_zh':
      return { image: word_zh, text: '传统文档' }
    case 'word_en':
      return { image: word_en, text: 'Document' }
    case 'image':
    default:
      return { image: image, text: lang == 'zh-CN' ? '图片' : 'Image' }
  }
}