import {Box, BoxProps, FlexProps} from "@chakra-ui/react";
import {PropsWithChildren} from "react";
import {keyframes} from "@emotion/react";

const fadeIn = keyframes`
    from {
        opacity: 0;
        transform: translateY(1rem);
        filter: blur(6px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
        filter: blur(0)
    }
`;

interface FadeInBoxProps extends PropsWithChildren<BoxProps & FlexProps> {
    isAnimationDisabled?: boolean;
}

export default function FadeInBox({children, isAnimationDisabled = false, ...props}: FadeInBoxProps) {
    return (
        <Box {...!isAnimationDisabled && {animation: `${fadeIn} ease-in-out 1s both`}}
             {...props}>
            {children}
        </Box>
    );
}