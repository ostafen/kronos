import {BoxProps, Text} from "@chakra-ui/react";
import {PropsWithChildren} from "react";

interface GradientTextProps extends BoxProps {
    text: string;
}

export default function GradientText({text, children, ...boxProps}: PropsWithChildren<GradientTextProps>) {
    return (
        <Text backgroundClip="text"
              as="span"
              color="transparent"
              width="100%"
              bgImage="linear-gradient(to right, #e66465, #9198e5)"
              {...boxProps}
        >
            {children}
            {text}
        </Text>
    )
}