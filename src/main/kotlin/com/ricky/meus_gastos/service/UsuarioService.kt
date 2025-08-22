package com.ricky.meus_gastos.service

import com.ricky.meus_gastos.models.Usuario

interface UsuarioService {
    fun findByEmail(email:String):Usuario
}