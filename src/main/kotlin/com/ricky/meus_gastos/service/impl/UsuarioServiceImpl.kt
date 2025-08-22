package com.ricky.meus_gastos.service.impl

import com.ricky.meus_gastos.exceptions.GenericException
import com.ricky.meus_gastos.models.Usuario
import com.ricky.meus_gastos.repository.UsuarioRepository
import com.ricky.meus_gastos.service.UsuarioService
import org.springframework.stereotype.Service

@Service
class UsuarioServiceImpl(private val usuarioRepository: UsuarioRepository) : UsuarioService {
    override fun findByEmail(email: String): Usuario =
        usuarioRepository.findByEmail(email) ?: throw GenericException("email.nao.encotrado")
}