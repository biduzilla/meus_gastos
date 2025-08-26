package com.ricky.meus_gastos.service.impl

import com.ricky.meus_gastos.dto.LoginDTO
import com.ricky.meus_gastos.dto.TokenDTO
import com.ricky.meus_gastos.exceptions.GenericException
import com.ricky.meus_gastos.models.Usuario
import com.ricky.meus_gastos.repository.UsuarioRepository
import com.ricky.meus_gastos.security.JwtService
import com.ricky.meus_gastos.service.UsuarioService
import org.springframework.http.HttpStatus
import org.springframework.security.authentication.AuthenticationManager
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken
import org.springframework.security.core.AuthenticationException
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional

@Service
class UsuarioServiceImpl(
    private val usuarioRepository: UsuarioRepository,
    private val passwordEncoder: PasswordEncoder,
    private val authenticationManager: AuthenticationManager,
    private val jwtService: JwtService
) : UsuarioService {
    override fun findByEmail(email: String): Usuario {
        return usuarioRepository.findByEmail(email) ?: throw GenericException("email.nao.encontrado")
    }

    override fun login(login: LoginDTO): TokenDTO {
        try {
            authenticationManager.authenticate(
                UsernamePasswordAuthenticationToken(
                    login.login,
                    login.senha
                )
            )

            val usuario = findByEmail(login.login)
            val token = jwtService.generateToken(usuario)
            return TokenDTO(
                token = token,
                idUser = usuario.idUsuario,
                nome = usuario.nome
            )
        } catch (e: AuthenticationException) {
            throw GenericException(
                msg = "error.login.invalido",
                httpStatus = HttpStatus.BAD_REQUEST
            )
        }
    }

    @Transactional
    override fun update(usuario: Usuario): Usuario {
        if (usuarioRepository.existsByEmail(usuario.email)) {
            throw GenericException(
                msg = "error.email.cadastrado",
                httpStatus = HttpStatus.BAD_REQUEST
            )
        }
        val user = findById(usuario.idUsuario).copy(
            nome = usuario.nome,
            email = usuario.email,
            senha = passwordEncoder.encode(usuario.senha)
        )

        return save(user)
    }


    @Transactional
    override fun save(usuario: Usuario): Usuario {
        if (usuarioRepository.existsByEmail(usuario.email)) {
            throw GenericException(
                msg = "error.email.cadastrado",
                httpStatus = HttpStatus.BAD_REQUEST
            )
        }
        usuario.senha = passwordEncoder.encode(usuario.senha)
        return usuarioRepository.save(usuario)
    }

    override fun deleteById(idUsuario: String) {
        usuarioRepository.deleteById(idUsuario)
    }

    override fun findById(idUsuario: String): Usuario {
        return usuarioRepository.findById(idUsuario).orElseThrow {
            throw GenericException(
                msg = "usuario.nao.encotrado",
                httpStatus = HttpStatus.NOT_FOUND
            )
        }
    }
}