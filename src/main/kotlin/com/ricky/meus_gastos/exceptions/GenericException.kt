package com.ricky.meus_gastos.exceptions

import org.springframework.http.HttpStatus

class GenericException(
    private val msg: String,
    val httpStatus: HttpStatus = HttpStatus.BAD_REQUEST
) : RuntimeException(msg) {
}