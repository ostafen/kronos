import {IconButton, IconButtonProps} from "@chakra-ui/react";
import {Link, LinkProps} from "react-router";
import {PropsWithChildren} from "react";

type IconButtonLinkProps = LinkProps & IconButtonProps;

export default function IconButtonLink({children, ...props}: PropsWithChildren<IconButtonLinkProps>) {
    return (
        <IconButton as={Link} {...props}>
            {children}
        </IconButton>
    );
}