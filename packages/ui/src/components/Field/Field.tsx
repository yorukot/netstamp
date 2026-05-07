import * as Label from "@radix-ui/react-label";
import type { ComponentPropsWithoutRef, ReactNode } from "react";
import { forwardRef, useId } from "react";
import styles from "./Field.module.css";

export type ControlVariant = "default" | "compact" | "bare";

function booleanAria(value: boolean | undefined) {
	return value ? true : undefined;
}

export interface InputProps extends ComponentPropsWithoutRef<"input"> {
	variant?: ControlVariant;
	invalid?: boolean;
	frameClassName?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(function Input({ variant = "default", invalid, frameClassName, className, "aria-invalid": ariaInvalidProp, ...props }, ref) {
	const ariaInvalid = invalid || ariaInvalidProp === true || ariaInvalidProp === "true";
	const classes = [styles.control, styles[`${variant}Control`], className].filter(Boolean).join(" ");

	if (variant === "bare") {
		return <input ref={ref} className={classes} aria-invalid={booleanAria(ariaInvalid)} {...props} />;
	}

	const frameClasses = ["ns-cut-frame", styles.controlFrame, styles[`${variant}Frame`], frameClassName].filter(Boolean).join(" ");

	return (
		<span className={frameClasses} data-invalid={Boolean(ariaInvalid)}>
			<input ref={ref} className={classes} aria-invalid={booleanAria(ariaInvalid)} {...props} />
		</span>
	);
});

export interface SelectProps extends ComponentPropsWithoutRef<"select"> {
	variant?: Exclude<ControlVariant, "bare">;
	invalid?: boolean;
	frameClassName?: string;
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(function Select(
	{ variant = "default", invalid, frameClassName, className, children, "aria-invalid": ariaInvalidProp, ...props },
	ref
) {
	const ariaInvalid = invalid || ariaInvalidProp === true || ariaInvalidProp === "true";
	const classes = [styles.control, styles.select, styles[`${variant}Control`], className].filter(Boolean).join(" ");
	const frameClasses = ["ns-cut-frame", styles.controlFrame, styles.selectFrame, styles[`${variant}Frame`], frameClassName].filter(Boolean).join(" ");

	return (
		<span className={frameClasses} data-invalid={Boolean(ariaInvalid)}>
			<select ref={ref} className={classes} aria-invalid={booleanAria(ariaInvalid)} {...props}>
				{children}
			</select>
		</span>
	);
});

export interface CheckboxProps extends Omit<ComponentPropsWithoutRef<"input">, "type"> {
	invalid?: boolean;
}

export const Checkbox = forwardRef<HTMLInputElement, CheckboxProps>(function Checkbox({ invalid, className, "aria-invalid": ariaInvalidProp, ...props }, ref) {
	const ariaInvalid = invalid || ariaInvalidProp === true || ariaInvalidProp === "true";
	const classes = [styles.checkbox, className].filter(Boolean).join(" ");

	return <input ref={ref} type="checkbox" className={classes} aria-invalid={booleanAria(ariaInvalid)} {...props} />;
});

export interface FieldShellProps {
	id: string;
	label: ReactNode;
	helper?: ReactNode;
	error?: ReactNode;
	children: ReactNode;
}

function FieldShell({ id, label, helper, error, children }: FieldShellProps) {
	return (
		<div className={styles.field}>
			<Label.Root className={styles.label} htmlFor={id}>
				{label}
			</Label.Root>
			{children}
			{error ? <span className={styles.error}>{error}</span> : null}
			{helper && !error ? <span className={styles.helper}>{helper}</span> : null}
		</div>
	);
}

export interface TextFieldProps extends ComponentPropsWithoutRef<"input"> {
	label: ReactNode;
	helper?: ReactNode;
	error?: ReactNode;
}

export function TextField({ label, helper, error, className, ...props }: TextFieldProps) {
	const generatedId = useId();
	const id = props.id || generatedId;

	return (
		<FieldShell id={id} label={label} helper={helper} error={error}>
			<Input id={id} className={className} invalid={Boolean(error)} {...props} />
		</FieldShell>
	);
}

export interface TextAreaFieldProps extends ComponentPropsWithoutRef<"textarea"> {
	label: ReactNode;
	helper?: ReactNode;
	error?: ReactNode;
}

export function TextAreaField({ label, helper, error, className, ...props }: TextAreaFieldProps) {
	const generatedId = useId();
	const id = props.id || generatedId;
	const classes = [styles.control, styles.area, className].filter(Boolean).join(" ");

	return (
		<FieldShell id={id} label={label} helper={helper} error={error}>
			<span className={["ns-cut-frame", styles.controlFrame].join(" ")} data-invalid={Boolean(error)}>
				<textarea id={id} className={classes} aria-invalid={Boolean(error)} {...props} />
			</span>
		</FieldShell>
	);
}

export interface SelectFieldProps extends ComponentPropsWithoutRef<"select"> {
	label: ReactNode;
	helper?: ReactNode;
	error?: ReactNode;
	options: Array<{ value: string; label: string }>;
}

export function SelectField({ label, helper, error, options, className, ...props }: SelectFieldProps) {
	const generatedId = useId();
	const id = props.id || generatedId;

	return (
		<FieldShell id={id} label={label} helper={helper} error={error}>
			<Select id={id} className={className} invalid={Boolean(error)} {...props}>
				{options.map(option => (
					<option key={option.value} value={option.value}>
						{option.label}
					</option>
				))}
			</Select>
		</FieldShell>
	);
}
