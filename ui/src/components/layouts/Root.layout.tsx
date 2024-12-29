import Navbar from '../organisms/Navbar/Navbar';
import {Outlet} from 'react-router';
import {Container} from "@chakra-ui/react";

export default function RootLayout() {
    return (
        <>
            <Navbar mb="10"/>
            <Container fluid>
                <Outlet/>
            </Container>
        </>
    );
}
