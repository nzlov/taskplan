import axios from 'axios';
import md5 from 'js-md5';

function genSign(user, token, url) {
  return md5(`${user}/api${url}${token}`);
}

class HttpUtil {
  static Get(url) {
    return axios.get(`http://192.168.1.200:9005/api${url}`);
  }
  static Post(url, data) {
    return axios.post(`http://192.168.1.200:9005/api${url}`, data);
  }
  static LGet(store, url) {
    return axios({
      method: 'get',
      url: `http://192.168.1.200:9005/api${url}`,
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
      url: `http://192.168.1.200:9005/api${url}`,
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
      url: `http://192.168.1.200:9005/api${url}`,
      headers: {
        'X-AppUser': store.user,
        'X-AppSign': genSign(store.user, store.token, url),
      },
    });
  }
}
export default HttpUtil;
