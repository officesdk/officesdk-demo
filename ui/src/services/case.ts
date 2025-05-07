import Request from '@/services/request';

export default {
  GetCases: () => {
    return Request(`/api/example/case`);
  },

  GetCasePreview: (fileId: string) => {
    return Request(`/showcase/${fileId}/page`)
  }
};
