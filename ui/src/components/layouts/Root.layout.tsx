import Navbar from '../organisms/Navbar/Navbar';
import { Outlet } from 'react-router';

export default function RootLayout() {
  return (
    <>
      <Navbar mb="10" />
      <Outlet />
    </>
  );
}
