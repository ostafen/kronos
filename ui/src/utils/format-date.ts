import dayjs from 'dayjs';

export default function formatDate(date: string) {
  return dayjs(date).format('YYYY/MM/DD HH:mm:ss');
}
