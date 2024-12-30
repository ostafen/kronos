import {Outlet} from 'react-router';
import {Container} from "@chakra-ui/react";

export default function RootLayout() {
    return (
        <Container fluid pb="4rem">
            <Outlet/>
        </Container>
    );
}
