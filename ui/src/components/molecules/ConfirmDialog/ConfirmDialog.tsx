import { Button } from "@/components/chakra/button";
import {
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogRoot,
  DialogTitle,
} from "@/components/chakra/dialog";
import { DialogActionTrigger } from "@chakra-ui/react";
import { ReactNode, useEffect, useState } from "react";
import { Subject } from "rxjs";

interface ConfirmDialogProps {
  title: string;
  content: ReactNode;
  isOpen: boolean;
}

export const dialogClose$ = new Subject<void>();
export const dialogConfirm$ = new Subject<void>();
export const dialogReset$ = new Subject<void>();

export default function ConfirmDialog(props: ConfirmDialogProps) {
  const [isConfirmed, setIsConfirmed] = useState(false);
  const { title, content, isOpen } = props;

  const handleConfirm = () => {
    setIsConfirmed(true);
    dialogConfirm$.next();
  };

  const handleCancel = () => {
    if (isConfirmed) return;
    dialogClose$.next();
  };

  useEffect(() => {
    const sub = dialogReset$.subscribe(() => setIsConfirmed(false));
    return () => sub.unsubscribe();
  }, []);

  return (
    <DialogRoot open={isOpen} onOpenChange={handleCancel} placement="center">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
        </DialogHeader>
        <DialogBody>{content}</DialogBody>
        <DialogFooter>
          <DialogActionTrigger asChild>
            <Button variant="outline">Cancel</Button>
          </DialogActionTrigger>
          <Button onClick={handleConfirm}>Confirm</Button>
        </DialogFooter>
        <DialogCloseTrigger />
      </DialogContent>
    </DialogRoot>
  );
}
