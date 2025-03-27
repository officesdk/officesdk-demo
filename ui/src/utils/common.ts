
import example_e from '../assets/example_e.png'
import example_w from '../assets/example_w.png'
import example_p from '../assets/example_p.png'
import example_i from '../resource/image.png'

// 获取对应case信息
export function getImage(type: string): { image: string, text: string } {
  switch (type) {
    case 'excel':
      return { image: example_e, text: '表格' }
    case 'ppt':
      return { image: example_p, text: '演示文档' }
    case 'word':
      return { image: example_w, text: '传统文档' }
    case 'image':
    default:
      return { image: example_i, text: '图片' }
  }
}