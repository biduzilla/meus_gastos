package com.ricky.meus_gastos

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication
import org.springframework.data.jpa.repository.config.EnableJpaAuditing

@SpringBootApplication
@EnableJpaAuditing(auditorAwareRef = "auditAwareImpl")
class MeusGastosApplication

fun main(args: Array<String>) {
	runApplication<MeusGastosApplication>(*args)
}
