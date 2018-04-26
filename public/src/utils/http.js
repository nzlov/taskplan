import axios from 'axios';
import md5 from 'js-md5';

function genSign(user, token, url) {
  return md5(`${user}/api${url}${token}`);
}

const publicurl = '';

class HttpUtil {
  static Get(url) {
    return axios.get(`${publicurl}/api${url}`);
  }
  static Post(url, data) {
    return axios.post(`${publicurl}/api${url}`, data);
  }
  static LGet(store, url) {
    return axios({
      method: 'get',
      url: `${publicurl}/api${url}`,
      headers: {
        'X-AppUser': store.user,
        'X-AppSign': genSign(store.user, store.token, url),
      },
    });
  }
  static LPost(store, url, data) {
    console.dir(data);
    return axios({
      method: 'post',
      url: `${publicurl}/api${url}`,
      headers: {
        'X-AppUser': store.user,
        'X-AppSign': genSign(store.user, store.token, url),
      },
      data,
    });
  }
  static LDel(store, url) {
    return axios({
      method: 'delete',
      url: `${publicurl}/api${url}`,
      headers: {
        'X-AppUser': store.user,
        'X-AppSign': genSign(store.user, store.token, url),
      },
    });
  }
}
export default HttpUtil;
