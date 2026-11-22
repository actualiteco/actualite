import { useFormik } from 'formik';
import { toast } from 'react-toastify';
import TextField from '@components/shared/TextField/TextField';
import Button from '@components/shared/Button/Button';
import { useAuthStore } from '@stores/useAuthStore';
import { parseApiError } from '@utils/error.utils';
import { isFilled, signinSchema } from './validation';
import styles from './AuthForms.module.scss';

interface SigninFormValues {
    email: string;
    password: string;
}

const SigninForm = () => {
    const login = useAuthStore((state) => state.login);
    const isLoading = useAuthStore((state) => state.isLoading);

    const formik = useFormik<SigninFormValues>({
        initialValues: { email: '', password: '' },
        validationSchema: signinSchema,
        onSubmit: async (values: SigninFormValues) => {
            try {
                await login(values.email.trim(), values.password);
                toast.success('С возвращением!');
            } catch (error) {
                const { status, message } = parseApiError(error);
                if (status === 401 || status === 400 || status === 403) {
                    toast.error('Неверный email или пароль');
                } else if (status === 429) {
                    toast.error('Слишком много попыток. Попробуйте позже');
                } else {
                    toast.error(message || 'Не удалось войти. Попробуйте ещё раз');
                }
            }
        },
    });

    const disabled = !isFilled({ ...formik.values });

    return (
        <form onSubmit={formik.handleSubmit} className={styles.form} noValidate>
            <h1 className={styles.title}>С возвращением</h1>
            <p className={styles.subtitle}>Продолжайте делиться мыслями</p>

            <TextField
                name="email"
                type="email"
                placeholder="Email"
                autoComplete="email"
                value={formik.values.email}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.errors.email}
                touched={formik.touched.email}
            />

            <TextField
                name="password"
                type="password"
                placeholder="Пароль"
                autoComplete="current-password"
                value={formik.values.password}
                onChange={formik.handleChange}
                onBlur={formik.handleBlur}
                error={formik.errors.password}
                touched={formik.touched.password}
            />

            <Button type="submit" fullWidth isLoading={isLoading} disabled={disabled}>
                Войти
            </Button>
        </form>
    );
};

export default SigninForm;
