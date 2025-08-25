package com.ricky.meus_gastos.controller

import com.ricky.meus_gastos.dto.LoginDTO
import com.ricky.meus_gastos.dto.TokenDTO
import com.ricky.meus_gastos.dto.UsuarioDTO
import com.ricky.meus_gastos.dto.toModel
import com.ricky.meus_gastos.models.toDTO
import com.ricky.meus_gastos.service.UsuarioService
import jakarta.validation.Valid
import org.springframework.web.bind.annotation.*

@RestController
@RequestMapping("/usuario")
class UsuarioController(
    private val usuarioService: UsuarioService
) {
    @GetMapping("/login")
    suspend fun login(@RequestBody loginDTO: LoginDTO): TokenDTO {
        return usuarioService.login(loginDTO)
    }

    @PostMapping("/save")
    suspend fun save(@RequestBody @Valid usuarioDTO: UsuarioDTO): UsuarioDTO {
        return usuarioService.save(usuarioDTO.toModel()).toDTO()
    }

    @PutMapping("/update")
    suspend fun update(@RequestBody @Valid usuarioDTO: UsuarioDTO): UsuarioDTO {
        return usuarioService.update(usuarioDTO.toModel()).toDTO()
    }

    @PutMapping("/delete/{id}")
    suspend fun deleteById(@PathVariable id: String) {
        usuarioService.deleteById(id)
    }
}