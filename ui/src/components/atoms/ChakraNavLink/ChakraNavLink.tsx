import {
  Link as ChakraLink,
  LinkProps as ChakraLinkProps,
} from '@chakra-ui/react';
import { PropsWithChildren } from 'react';
import { LinkProps, NavLink } from 'react-router';

type ChakraNavLinkProps = LinkProps & ChakraLinkProps;

export default function ChakraNavLink({
  children,
  ...props
}: PropsWithChildren<ChakraNavLinkProps>) {
  return (
    <ChakraLink as={NavLink} {...props}>
      {children}
    </ChakraLink>
  );
}
