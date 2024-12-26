import { Container, Heading } from '@chakra-ui/react';

import {
  BreadcrumbRoot,
  BreadcrumbCurrentLink,
} from '@/components/chakra/breadcrumb';
import { FiChevronRight } from 'react-icons/fi';
import ChakraBreadcrumbLink from '@/components/atoms/ChakraBreadcrumbLink/ChakraBreadcrumbLink';

export default function NewSchedulePage() {
  return (
    <Container>
      <BreadcrumbRoot separator={<FiChevronRight />} variant="underline">
        <ChakraBreadcrumbLink to="/">Home</ChakraBreadcrumbLink>
        <BreadcrumbCurrentLink>New Schedule</BreadcrumbCurrentLink>
      </BreadcrumbRoot>
      <Heading mt="3">New Schedule</Heading>
    </Container>
  );
}
