import * as Yup from 'yup';
import {
    DISPLAY_NAME_MAX,
    STATUS_MAX,
    USERNAME_MAX,
    USERNAME_MIN,
    USERNAME_REGEX,
    displayNameField,
    statusField,
    usernameField,
} from '@utils/validators.utils';

export { DISPLAY_NAME_MAX, STATUS_MAX };

export const profileSchema = Yup.object().shape({
    displayName: displayNameField(),
    username: usernameField(),
    status: statusField(),
});

export const isUsernameFormatValid = (username: string): boolean =>
    username.length >= USERNAME_MIN &&
    username.length <= USERNAME_MAX &&
    USERNAME_REGEX.test(username);
