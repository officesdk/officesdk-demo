import { extend } from 'umi-request';

const prefix = process.env.NODE_ENV === 'development' ? '/api' : process.env.BASE ? process.env.BASE.endsWith('/') ? process.env.BASE.substring(0, process.env.BASE.length - 1) : process.env.BASE : '';

const Request = extend({
  prefix: prefix,
});

export default Request;
