package com.ricky.meus_gastos.utils

import org.springframework.context.MessageSource
import org.springframework.stereotype.Component
import java.util.*

@Component
class I18n(private val messageSource: MessageSource) {

    fun getMessage(message: String): String {
        return messageSource.getMessage(message, null, Locale.getDefault())
    }

    fun getMessage(message: String, vararg args: Any): String {
        return messageSource.getMessage(message, args, Locale.getDefault())
    }

    fun getMessage(message: String, locale: Locale): String {
        return messageSource.getMessage(message, null, locale)
    }

    fun getMessage(message: String, vararg args: Any, locale: Locale): String {
        return messageSource.getMessage(message, args, locale)
    }
}