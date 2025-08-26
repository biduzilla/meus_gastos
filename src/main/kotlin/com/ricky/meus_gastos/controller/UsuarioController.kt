package com.ricky.meus_gastos.controller

import com.ricky.meus_gastos.dto.*
import com.ricky.meus_gastos.models.toDTO
import com.ricky.meus_gastos.service.UsuarioService
import jakarta.validation.Valid
import org.springframework.web.bind.annotation.*

@RestController
@RequestMapping("/usuario")
class UsuarioController(
    private val usuarioService: UsuarioService
) {
    @GetMapping("/{idUsuario}")
    fun findById(@PathVariable idUsuario:String): UsuarioDTO{
        return usuarioService.findById(idUsuario).toDTO()
    }

    @PostMapping("/login")
    fun login(@RequestBody @Valid loginDTO: LoginDTO): TokenDTO {
        return usuarioService.login(loginDTO)
    }

    @PostMapping("/save")
    fun save(@RequestBody @Valid usuarioSaveDTO: UsuarioSaveDTO): UsuarioDTO {
        return usuarioService.save(usuarioSaveDTO.toModel()).toDTO()
    }

    @PutMapping("/update")
    fun update(@RequestBody @Valid usuarioSaveDTO: UsuarioSaveDTO): UsuarioDTO {
        return usuarioService.update(usuarioSaveDTO.toModel()).toDTO()
    }

    @DeleteMapping("/delete/{id}")
    fun deleteById(@PathVariable id: String) {
        usuarioService.deleteById(id)
    }
}