class DateUtil {

  // 1451718245 to '2016-01-02 15:04:05'
  static formatDate(d) {
    if (d <= 0) {
      return '';
    }
    const now = new Date(d * 1000);
    const year = now.getFullYear();
    let month = now.getMonth() + 1;
    let date = now.getDate();
    let hour = now.getHours();
    let minute = now.getMinutes();
    let second = now.getSeconds();
    if (month <= 9) {
      month = `0${month}`;
    }
    if (date <= 9) {
      date = `0${date}`;
    }
    if (hour <= 9) {
      hour = `0${hour}`;
    }
    if (minute <= 9) {
      minute = `0${minute}`;
    }
    if (second <= 9) {
      second = `0${second}`;
    }
    return `${year}-${month}-${date} ${hour}:${minute}:${second}`;
  }
  // Date to '2016-01-02'
  static formatDate1(now) {
    if (now <= 0) {
      return '';
    }
    const year = now.getFullYear();
    let month = now.getMonth() + 1;
    let date = now.getDate();
    if (month <= 9) {
      month = `0${month}`;
    }
    if (date <= 9) {
      date = `0${date}`;
    }
    return `${year}-${month}-${date}`;
  }
  // Date to '15:04'
  static formatTime(now) {
    if (now <= 0) {
      return '';
    }
    let hour = now.getHours();
    let minute = now.getMinutes();
    if (hour <= 9) {
      hour = `0${hour}`;
    }
    if (minute <= 9) {
      minute = `0${minute}`;
    }
    return `${hour}:${minute}`;
  }
  // '2016-01-02 15:04:05' to 1451718245
  static formatTimestamp(now) {
    if (!now) {
      return 0;
    }
    return Date.parse(now) / 1000;
  }
}

export default DateUtil;
