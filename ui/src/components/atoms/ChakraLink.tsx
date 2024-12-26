import { LinkProps as ChakraLinkProps, Link } from '@chakra-ui/react';
import { PropsWithChildren } from 'react';
import { Link as ReactRouterLink } from 'react-router';
import { LinkProps as ReactRouterLinkProps } from 'react-router';

type LinkProps = ReactRouterLinkProps & ChakraLinkProps;

export default function ChakraLink({
  children,
  ...props
}: PropsWithChildren<LinkProps>) {
  return (
    <Link as={ReactRouterLink} {...props}>
      {children}
    </Link>
  );
}
