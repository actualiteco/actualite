import { useFormik } from 'formik';
import { toast } from 'react-toastify';
import TextField from '@components/shared/TextField/TextField';
import Button from '@components/shared/Button/Button';
import { useAuthStore } from '@stores/useAuthStore';
import { getApiErrorMessage, parseApiError } from '@utils/error.utils';
import { isFilled, signupSchema } from './validation';
import styles from './AuthForms.module.scss';

interface SignupFormValues {
    username: string;
    email: string;
    password: string;
}

const EMAIL_CODE = /EMAIL/;
const USERNAME_CODE = /USERNAME|NICKNAME/;
const TAKEN_CODE = /ALREADY|USED|EXIST|TAKEN|DUPLICATE|CONFLICT/;

const TAKEN_HINT = /taken|exist|already|used|занят|существ|использ/;

const SignupForm = () => {
    const register = useAuthStore((state) => state.register);
    const isLoading = useAuthStore((state) => state.isLoading);

    const formik = useFormik<SignupFormValues>({
        initialValues: { username: '', email: '', password: '' },
        validationSchema: signupSchema,
        onSubmit: async (values: SignupFormValues) => {
            try {
                await register(values.email.trim(), values.username, values.password);
                toast.success('Добро пожаловать в actualite!');
            } catch (error) {
                const { status, code, message } = parseApiError(error);
                const errorCode = (code ?? '').toUpperCase();
                const text = message.toLowerCase();

                if (status === 429) {
                    toast.error('Слишком много попыток. Попробуйте позже');
                } else if (TAKEN_CODE.test(errorCode) && EMAIL_CODE.test(errorCode)) {
                    formik.setFieldError('email', 'Email уже используется');
                } else if (TAKEN_CODE.test(errorCode) && USERNAME_CODE.test(errorCode)) {
                    formik.setFieldError('username', 'Никнейм занят');
                } else if (/username|никнейм|имя пользователя/.test(text) && TAKEN_HINT.test(text)) {
                    formik.setFieldError('username', 'Никнейм занят');
                } else if (/email|почт/.test(text) && TAKEN_HINT.test(text)) {
                    formik.setFieldError('email', 'Email уже используется');
                } else if (status === 409) {
                    formik.setFieldError('email', 'Email или никнейм уже используются');
                } else {
                    toast.error(
                        getApiErrorMessage(error, 'Не удалось зарегистрироваться. Попробуйте ещё раз'),
                    );
                }
            }
        },
    });

    const shown = (name: keyof SignupFormValues) =>
        Boolean(formik.touched[name]) || formik.values[name] !== '';

    const disabled = !isFilled({ ...formik.values }) || !formik.isValid;

    return (
        <form onSubmit={formik.handleSubmit} className={styles.form} noValidate>
            <h1 className={styles.title}>Создайте аккаунт</h1>
            <p className={styles.subtitle}>Продолжайте делиться мыслями</p>

            <TextField
                name="username"
                type="text"
                placeholder="Имя пользователя"
                autoComplete="username"
                value={formik.values.username}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.errors.username}
                touched={shown('username')}
            />

            <TextField
                name="email"
                type="email"
                placeholder="Email"
                autoComplete="email"
                value={formik.values.email}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.errors.email}
                touched={shown('email')}
            />

            <TextField
                name="password"
                type="password"
                placeholder="Пароль"
                autoComplete="new-password"
                value={formik.values.password}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.errors.password}
                touched={shown('password')}
            />

            <Button type="submit" fullWidth isLoading={isLoading} disabled={disabled}>
                Зарегистрироваться
            </Button>
        </form>
    );
};

export default SignupForm;
