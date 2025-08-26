package com.ricky.meus_gastos.models

import jakarta.persistence.Column
import jakarta.persistence.EntityListeners
import jakarta.persistence.MappedSuperclass
import org.springframework.data.annotation.CreatedBy
import org.springframework.data.annotation.CreatedDate
import org.springframework.data.annotation.LastModifiedBy
import org.springframework.data.annotation.LastModifiedDate
import org.springframework.data.jpa.domain.support.AuditingEntityListener
import java.io.Serializable
import java.time.LocalDateTime

@MappedSuperclass
@EntityListeners(AuditingEntityListener::class)
abstract class BaseModel(
    @Column(updatable = false, name = "CREATED_AT")
    @field:CreatedDate
    var createdAt: LocalDateTime? = null,

    @field:CreatedBy
    @Column(updatable = false, name = "CREATED_BY")
    var createdBy: String? = null,

    @field:LastModifiedDate
    @Column(insertable = false, name = "UPDATED_AT")
    var updatedAt: LocalDateTime? = null,

    @field:LastModifiedBy
    @Column(insertable = false, name = "UPDATED_BY")
    var updatedBy: String? = null,

    @Column(name = "DELETED")
    var deleted: Boolean = false,
) : Serializable
