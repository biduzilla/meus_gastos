package com.ricky.meus_gastos.config

import org.springframework.context.MessageSource
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.context.support.ResourceBundleMessageSource

@Configuration
class MessageConfiguration {

    @Bean
    fun messageSource(): MessageSource {
        val source = ResourceBundleMessageSource()
        source.setBasename("lang/messages")
        source.setUseCodeAsDefaultMessage(true)

        return source
    }
}