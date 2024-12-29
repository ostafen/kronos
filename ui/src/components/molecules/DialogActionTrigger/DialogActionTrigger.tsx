import {IconButtonProps} from "@chakra-ui/react";
import {PropsWithChildren} from "react";
import {Button} from "@/components/chakra/button.tsx";
import {dialogClose$, dialogConfirm$, dialogOpen$, dialogReset$} from "@/utils/modal.ts";
import DialogData from "@/model/dialog-data.ts";

interface DeleteScheduleTriggerProps extends IconButtonProps {
    onConfirm: () => Promise<void>;
    onSuccess?: () => void;
    dialogData: Omit<DialogData, 'isConfirmed'>;
}

export default function DialogActionTrigger(props: PropsWithChildren<DeleteScheduleTriggerProps>) {
    const {children, onSuccess, onConfirm, dialogData, ...iconButtonProps} = props;

    const initFlow = async () => {
        dialogOpen$.sink.next({
            ...dialogData,
            isConfirmed: false
        });

        try {
            await new Promise((resolve, reject) => {
                const confirmSub = dialogConfirm$.source$
                    .subscribe(() => {
                        resolve({status: "confirmed"});
                        confirmSub.unsubscribe();
                    });

                const cancelSub = dialogClose$.source$
                    .subscribe(() => {
                        reject({status: "canceled"});
                        cancelSub.unsubscribe();
                    });
            });
            await onConfirm();
            dialogReset$.sink.next();
            onSuccess?.();
        } catch (error) {
            if (error instanceof Error) {
                console.error(error.message);
            }
        } finally {
            dialogClose$.sink.next();
        }
    };

    return (
        <Button
            onClick={initFlow}
            {...iconButtonProps}
        >
            {children}
        </Button>
    );
}