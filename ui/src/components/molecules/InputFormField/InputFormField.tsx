import {Field} from "@/components/chakra/field";
import FormField from "@/model/form-field";
import {Input, Textarea} from "@chakra-ui/react";
import {Controller, useFormContext} from "react-hook-form";
import {Switch} from "@/components/chakra/switch.tsx";

export default function InputFormField(props: FormField) {
    const {name, label, isRequired, type} = props;
    const {formState: {errors}} = useFormContext();
    const isInvalid = !!errors[name];
    const errorText = errors[name]?.message;

    return (
        <Field
            {...(typeof errorText === 'string' && {errorText})}
            {...(type !== 'checkbox' && {label})}
            required={isRequired}
            invalid={isInvalid}>
            <FieldComponent {...props} />
        </Field>
    )
}

const FieldComponent = (props: FormField) => {
    const {type, name, label, isRequired} = props;
    const {register} = useFormContext();

    const fieldProps = register(name, {
        ...(isRequired && {required: `${label} is required`})
    })

    switch (type) {
        case "textarea":
            return (
                <Textarea {...fieldProps}/>
            );

        case "input-text":
            return (
                <Input {...fieldProps} />
            );

        case "input-url":
            return (
                <Input type="url" {...fieldProps} />
            );

        case "input-cron":
            return (
                <Input {...fieldProps} />
            );

        case "input-datetime-local":
            return (
                <Input type="datetime-local" {...fieldProps} />
            );

        case "checkbox":
            return (
                <Controller name={name} render={({field}) => (
                    <Switch
                        name={field.name}
                        checked={field.value}
                        onCheckedChange={({checked}) => field.onChange(checked)}
                        inputProps={{onBlur: field.onBlur}}
                    >
                        {label}
                    </Switch>
                )}/>
            );
    }

    return null;
}
