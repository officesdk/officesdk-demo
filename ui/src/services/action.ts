import Request from '@/services/request';

export default {
  GetFiles: () => {
    return Request(`/showcase/files`)
  },
  DeleteFiles: (fileId: string) => {
    return Request(`/showcase/file/${fileId}`, {
      method: "DELETE"
    })
  },
  UploadFile: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return Request(`/showcase/file`, {
      method: 'POST',
      data: formData
    });
  }
};
