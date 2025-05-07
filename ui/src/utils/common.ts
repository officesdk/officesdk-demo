
import word_zh from '../assets/word_zh.png'
import word_en from '../assets/word_en.png'
import excel_zh from '../assets/excel_zh.png'
import excel_en from '../assets/excel_en.png'
import ppt_zh from '../assets/ppt_zh.png'
import ppt_en from '../assets/ppt_en.png'
import pdf_zh from '../assets/pdf_zh.png'
import pdf_en from '../assets/pdf_en.png'
import { FileType } from '@officesdk/web';

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
    case 'pdf_zh':
      return { image: pdf_zh, text: 'PDF' }
    case 'pdf_en':
      return { image: pdf_en, text: 'PDF' }
    default:
      return { image: "", text: lang == 'zh-CN' ? '未知类型' : 'Unknown' }
  }
}

export function getFileTypeFromExt(ext: string) {
  switch (ext) {
    case ".doc":
    case ".docx":
      return FileType.Document
    case ".ppt":
    case ".pptx":
      return FileType.Presentation
    case ".xls":
    case ".xlsx":
      return FileType.Spreadsheet
    case ".pdf":
      return FileType.Pdf
    default:
      return 0
  }
}
