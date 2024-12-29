import { MutableRefObject } from "react";

type ComponentWithRef<TProps, TRef> = TProps & {
  ref: MutableRefObject<TRef>;
};

export default ComponentWithRef;
