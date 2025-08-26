package com.ricky.meus_gastos.config

import com.ricky.meus_gastos.exceptions.GenericException
import com.ricky.meus_gastos.security.JwtService
import jakarta.servlet.FilterChain
import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import org.springframework.http.HttpStatus
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken
import org.springframework.security.core.context.SecurityContextHolder
import org.springframework.security.core.userdetails.UserDetailsService
import org.springframework.security.web.authentication.WebAuthenticationDetailsSource
import org.springframework.stereotype.Component
import org.springframework.web.filter.OncePerRequestFilter
import org.springframework.web.servlet.HandlerExceptionResolver

@Component
class JwtAuthenticationFilter(
    private val jwtService: JwtService,
    private val userDetailsService: UserDetailsService,
    private val handlerExceptionResolver: HandlerExceptionResolver
) : OncePerRequestFilter() {

    override fun doFilterInternal(
        request: HttpServletRequest,
        response: HttpServletResponse,
        filterChain: FilterChain
    ) {
        val authorization = request.getHeader("Authorization")

        if (authorization != null && authorization.startsWith("Bearer")) {
            try {
                val token = authorization.split(" ".toRegex()).dropLastWhile { it.isEmpty() }.toTypedArray()[1]
                if (jwtService.isTokenValid(token)) {
                    val loginUser: String = jwtService.extractUsername(token)
                    val user = userDetailsService.loadUserByUsername(loginUser)
                    val usuario = UsernamePasswordAuthenticationToken(user, null, user.authorities)
                    usuario.details = WebAuthenticationDetailsSource().buildDetails(request)
                    SecurityContextHolder.getContext().authentication = usuario
                } else {
                    val e = GenericException(
                        msg = "token.invalido",
                        httpStatus = HttpStatus.FORBIDDEN,
                    )
                    handlerExceptionResolver.resolveException(request, response, null, e)
                    return
                }
            } catch (e: Exception) {
                handlerExceptionResolver.resolveException(request, response, null, e)
                return
            }
        }
        filterChain.doFilter(request, response)
    }

//    override fun doFilterInternal(
//        request: HttpServletRequest,
//        response: HttpServletResponse,
//        filterChain: FilterChain
//    ) {
//        val authHeader: String? = request.getHeader("Authorization")
//
//        if (authHeader.isNullOrBlank() || !authHeader.startsWith("Bearer ")) {
//            filterChain.doFilter(request, response)
//            return
//        }
//
//        try {
//            val jwt = authHeader.substring(7)
//            val userEmail = jwtService.extractUsername(jwt)
//            val auth = SecurityContextHolder.getContext().authentication
//
//            if (!userEmail.isNullOrBlank() && auth == null) {
//                val userDetails = userDetailsService.loadUserByUsername(userEmail)
//
//                if (jwtService.isTokenValid(jwt, userDetails)) {
//                    val authToken = UsernamePasswordAuthenticationToken(
//                        userDetails,
//                        null,
//                        userDetails.authorities
//                    )
//                    authToken.details = WebAuthenticationDetailsSource().buildDetails(request)
//                    SecurityContextHolder.getContext().authentication = authToken
//                }
//            }
//
//            filterChain.doFilter(request, response)
//        } catch (e: Exception) {
//            handlerExceptionResolver.resolveException(request, response, null, e)
//        }
//
//    }
}