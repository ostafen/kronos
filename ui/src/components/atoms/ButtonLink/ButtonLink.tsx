import { ButtonProps } from '@/components/chakra/button';
import { Button } from '@chakra-ui/react';
import { PropsWithChildren } from 'react';
import { Link, LinkProps } from 'react-router';

type ButtonLinkProps = ButtonProps & LinkProps;

export default function ButtonLink({
  children,
  ...props
}: PropsWithChildren<ButtonLinkProps>) {
  return (
    <Button as={Link} {...props}>
      {' '}
      {children}
    </Button>
  );
}
