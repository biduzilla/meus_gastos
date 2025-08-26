package com.ricky.meus_gastos.config

import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration
import org.springframework.security.authentication.AuthenticationProvider
import org.springframework.security.config.annotation.web.builders.HttpSecurity
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity
import org.springframework.security.config.http.SessionCreationPolicy
import org.springframework.security.web.SecurityFilterChain
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter


@Configuration
@EnableWebSecurity
class SecurityConfiguration(
    private val authenticationProvider: AuthenticationProvider,
    private val jwtAuthenticationFilter: JwtAuthenticationFilter
) {
    @Bean
    @Throws(Exception::class)
    fun securityFilterChain(http: HttpSecurity): SecurityFilterChain {
        http.run {
            csrf { it.disable() }
            headers {
                it.frameOptions { frameOptions -> frameOptions.disable() } // ← IMPORTANTE!
            }
            authorizeHttpRequests {
                it.requestMatchers("/usuario/save").permitAll()
                it.requestMatchers("/usuario/login").permitAll()
                it.requestMatchers("/h2-console/**").permitAll()
                it.requestMatchers("/h2-console/").permitAll()
                it.requestMatchers("/h2-console").permitAll()
                    .anyRequest().authenticated()
            }
            sessionManagement {
                it.sessionCreationPolicy(SessionCreationPolicy.STATELESS)
            }
            authenticationProvider(authenticationProvider)
            addFilterBefore(jwtAuthenticationFilter, UsernamePasswordAuthenticationFilter::class.java)
        }

        return http.build()
    }
}