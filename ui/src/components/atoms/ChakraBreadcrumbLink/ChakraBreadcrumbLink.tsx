import { BreadcrumbLink, BreadcrumbLinkProps } from '@chakra-ui/react';
import { PropsWithChildren } from 'react';
import { LinkProps, Link } from 'react-router';

type ChakraBreadcrumbLinkProps = LinkProps & BreadcrumbLinkProps;

export default function ChakraBreadcrumbLink({
  children,
  ...props
}: PropsWithChildren<ChakraBreadcrumbLinkProps>) {
  return (
    <BreadcrumbLink as={Link} {...props}>
      {children}
    </BreadcrumbLink>
  );
}
