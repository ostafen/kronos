import {ReactNode} from "react";

export default interface DialogData {
    title: string;
    content: ReactNode;
    isConfirmed: boolean;
}
