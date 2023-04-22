import { loadTranslations, locales } from '$lib/translations';

function userLangCode() {
    if (typeof window == 'undefined') {
        return 'en'
    }
    const userLang = navigator.language || navigator.userLanguage;
    for (let i=0;i<locales.get().length;i++) {
        console.log(userLang, locales.get()[i]);
        if (userLang.includes(locales.get()[i])) {
            return locales.get()[i];
        }
    }
    return 'en';
}
export const load = async () => {
    await loadTranslations(userLangCode());
    return {};
}