import en from '@/utils/i18n/en.json'
import fr from '@/utils/i18n/fr.json'
import ko from '@/utils/i18n/ko.json'
import ru from '@/utils/i18n/ru.json'
import es from '@/utils/i18n/es.json'
import { useAppStore } from '@/store/app'

const trans = {
  'en': en,
  'fr': fr,
  'ko': ko,
  'ru': ru,
  'es': es,
}
export function T(key, params, num = 0) {
  const appStore = useAppStore()
  const lang = appStore.setting.lang
  const tran = trans[lang]?.[key]
  if (!tran) {
    return key
  }
  const msg = num > 1 ? (tran.Other ? tran.Other : tran.One) : tran.One
  // msg format: {name} is name
  // params format: {name: 'john'}
  // Replace placeholders
  return msg.replace(/{(\w+)}/g, function (match, key) {
    return params[key] || match
  })
}
