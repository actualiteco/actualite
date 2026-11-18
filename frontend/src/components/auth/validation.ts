import * as Yup from 'yup';
import { usernameField } from '@utils/validators.utils';

export { USERNAME_REGEX, USERNAME_MIN, USERNAME_MAX } from '@utils/validators.utils';

export const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/;
export const PRINTABLE_ASCII_REGEX = /^[\x20-\x7E]+$/;

export const PASSWORD_MIN = 8;
export const PASSWORD_MAX = 64;

const emailField = Yup.string()
    .trim()
    .required('Введите email')
    .matches(EMAIL_REGEX, 'Некорректный email');

export const signinSchema = Yup.object().shape({
    email: Yup.string().trim().required('Введите email'),
    password: Yup.string().required('Введите пароль'),
});

export const signupSchema = Yup.object().shape({
    username: usernameField(),
    email: emailField,
    password: Yup.string()
        .required('Введите пароль')
        .min(PASSWORD_MIN, `Пароль: от ${PASSWORD_MIN} до ${PASSWORD_MAX} символов`)
        .max(PASSWORD_MAX, `Пароль: от ${PASSWORD_MIN} до ${PASSWORD_MAX} символов`)
        .matches(PRINTABLE_ASCII_REGEX, 'Пароль может содержать только печатные символы ASCII'),
});

export const isFilled = (values: Record<string, string>): boolean =>
    Object.values(values).every((value) => value.trim() !== '');
