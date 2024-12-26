import { Box, BoxProps, Flex, HStack } from '@chakra-ui/react';
import ButtonLink from '@/components/atoms/ButtonLink/ButtonLink';
import ChakraLink from '@/components/atoms/ChakraLink';

const links: { name: string; path: string }[] = [];

export default function Navbar(props: BoxProps) {
  return (
    <Box px={4} py={2} borderBottom="1px solid rgba(0,0,0,.1)" {...props}>
      <Flex
        as="header"
        h={16}
        alignItems="center"
        justifyContent="space-between"
      >
        <ChakraLink to="/" fontSize="xl" fontWeight="bold" aria-label="Home">
          Kronos
        </ChakraLink>

        <nav aria-label="Primary Navigation">
          <HStack as="ul">
            {links.map((link) => (
              <li key={link.name}>
                <ButtonLink to={link.path}>{link.name}</ButtonLink>
              </li>
            ))}
          </HStack>
        </nav>
      </Flex>
    </Box>
  );
}
